#!/usr/bin/env python3
"""Read-only, real-Chromium acceptance of Wingless's self-contained ICE atlas.

This script exercises actual DOM input, selection, clicks, related-symbol navigation
and source-link popup behavior. The source popup is fulfilled locally: no network
traffic, research execution, acceptance mutation or external authorization occurs.
"""
import argparse
import json
import os
from pathlib import Path
from urllib.parse import urlparse

from playwright.sync_api import sync_playwright


def require(condition, message):
    if not condition:
        raise AssertionError(message)


def browser_path():
    override = os.environ.get("WINGLESS_CHROMIUM_EXE")
    paths = [Path(override)] if override else [
        Path(os.environ.get("PROGRAMFILES", "C:/Program Files")) / "Google/Chrome/Application/chrome.exe",
        Path(os.environ.get("PROGRAMFILES(X86)", "C:/Program Files (x86)")) / "Google/Chrome/Application/chrome.exe",
    ]
    for candidate in paths:
        if candidate.is_file():
            return candidate
    raise RuntimeError("Chrome (Chromium engine) not installed on this Windows runner")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--atlas", required=True, type=Path)
    parser.add_argument("--expected-source-commit", required=True)
    parser.add_argument("--receipt", required=True, type=Path)
    args = parser.parse_args()
    require(args.atlas.is_file(), "ATLAS_MISSING")
    require(len(args.expected_source_commit) == 40, "EXPECTED_SOURCE_SHA_NOT_FULL")
    receipt = {"schema": "wingless.atlas.chromium-acceptance.v1",
               "status": "FAIL", "atlas_path": str(args.atlas),
               "expected_source_commit": args.expected_source_commit,
               "checks": []}
    args.receipt.parent.mkdir(parents=True, exist_ok=True)
    page = None
    try:
        with sync_playwright() as playwright:
            browser = playwright.chromium.launch(
                executable_path=str(browser_path()), headless=True,
                args=["--no-first-run", "--no-default-browser-check"])
            receipt["browser_version"] = browser.version
            context = browser.new_context(accept_downloads=False)
            external_requests = []
            source_requests = []
            page_errors = []

            def intercept(route):
                url = route.request.url
                if url.startswith("https://github.com/spoonman136668-ai/Wingless/blob/"):
                    source_requests.append(url)
                    route.fulfill(status=200, content_type="text/html",
                                  body="<html><title>Offline source target probe</title>"
                                       "<body>Offline source target probe</body></html>")
                else:
                    external_requests.append(url)
                    route.abort()

            context.route("http://**/*", intercept)
            context.route("https://**/*", intercept)
            page = context.new_page()
            page.on("pageerror", lambda error: page_errors.append(str(error)))
            page.goto(args.atlas.resolve().as_uri(), wait_until="load")
            data = page.evaluate(
                "() => JSON.parse(document.getElementById('atlas').textContent)")
            nodes = data["nodes"]
            node_by_id = {node["id"]: node for node in nodes}
            require(data["source_commit"] == args.expected_source_commit,
                    "FROZEN_SOURCE_COMMIT_MISMATCH")
            require(data["authority"] == "advisory-navigation-only",
                    "ATLAS_AUTHORITY_BOUNDARY_DRIFT")
            require(len(nodes) == 352, "SOURCE_SYMBOL_COUNT_DRIFT")
            require(data["coverage"]["source_reviewed_labels"] == len(nodes),
                    "SOURCE_REVIEW_COVERAGE_DRIFT")
            require(all(n["evidence_basis"] == "source-reviewed" for n in nodes),
                    "UNREVIEWED_SYMBOL_ADMITTED")
            for key, expected in (
                ("nodesCount", len(nodes)),
                ("reviewedCount", len(nodes)),
                ("linksCount", len(data["relationships"])),
                ("sha", args.expected_source_commit[:12])):
                require(page.locator("#" + key).inner_text().strip() == str(expected),
                        "VISIBLE_STAT_MISMATCH:" + key)
            receipt["checks"].append("offline_page_and_statistics")

            # Exercise real search input and assert independently predicted result count.
            origin_edge = next((e for e in data["relationships"]
                                if e.get("from") in node_by_id
                                and e.get("to") in node_by_id
                                and e.get("from") != e.get("to")), None)
            require(origin_edge is not None, "NO_NAVIGABLE_SYMBOL_RELATIONSHIP")
            origin = node_by_id[origin_edge["from"]]
            target = node_by_id[origin_edge["to"]]
            needle = origin["label"].lower()
            expected_hits = sum(
                any(needle in str(field).lower() for field in (
                    n["label"], n["description"], n["file"], n["name"],
                    n.get("receiver", ""), *n.get("syntax_calls", [])))
                for n in nodes)
            page.locator("#q").fill(origin["label"])
            require(page.locator("#results button.symbol").count() == expected_hits
                    and expected_hits > 0, "SEARCH_DID_NOT_MATCH_MODEL")
            page.locator("#q").fill("")
            receipt["checks"].append("live_search_and_clear")

            page.locator("#basis").select_option("lexical-provisional")
            require(page.locator("#results button.symbol").count() == 0,
                    "PROVISIONAL_FILTER_NOT_EMPTY")
            page.locator("#basis").select_option("source-reviewed")
            require(page.locator("#results button.symbol").count() == len(nodes),
                    "SOURCE_REVIEW_FILTER_INCORRECT")
            page.locator("#basis").select_option("")
            kind = origin["kind"]
            page.locator("#kind").select_option(kind)
            require(page.locator("#results button.symbol").count() ==
                    sum(n["kind"] == kind for n in nodes),
                    "SYMBOL_KIND_FILTER_INCORRECT")
            page.locator("#kind").select_option("")
            receipt["checks"].append("evidence_and_kind_filters")

            # Click an actual list item, then follow a mapped relationship.
            origin_index = next(i for i, n in enumerate(nodes)
                                if n["id"] == origin["id"])
            page.locator("#results button.symbol").nth(origin_index).click()
            require(page.locator("#detail h3").first.inner_text() == origin["label"],
                    "SYMBOL_SELECTION_FAILED")
            require(page.locator("#results button.symbol").nth(
                origin_index).get_attribute("aria-current") == "true",
                "SYMBOL_SELECTION_NOT_VISIBLE")
            related_buttons = page.locator("#detail .related button")
            matching = [i for i, name in enumerate(related_buttons.all_inner_texts())
                        if name.strip() == "→ " + target["label"]]
            require(bool(matching), "RELATIONSHIP_LINK_NOT_RENDERED")
            related_buttons.nth(matching[0]).click()
            require(page.locator("#detail h3").first.inner_text() == target["label"],
                    "RELATIONSHIP_FOLLOW_FAILED")
            receipt["checks"].append("symbol_selection_and_relationship_navigation")

            # Source links must target frozen source, and open a real popup.
            source_link = page.locator("#detail a").filter(
                has_text="Open exact source definition")
            href = source_link.get_attribute("href")
            require(href == target["source_url"], "SOURCE_HREF_MISMATCH")
            parsed = urlparse(href)
            require(parsed.scheme == "https" and parsed.netloc == "github.com",
                    "SOURCE_URL_NOT_GITHUB_HTTPS")
            require(parsed.path.startswith(
                "/spoonman136668-ai/Wingless/blob/" +
                args.expected_source_commit + "/"), "SOURCE_URL_NOT_FROZEN")
            require(parsed.fragment == "L" + str(target["start_line"]) +
                    "-L" + str(target["end_line"]), "SOURCE_LINE_ANCHOR_DRIFT")
            require(source_link.get_attribute("target") == "_blank",
                    "SOURCE_LINK_NOT_NEW_TAB")
            require("noopener" in (source_link.get_attribute("rel") or ""),
                    "SOURCE_LINK_UNSAFE_REL")
            with page.expect_popup() as popup_event:
                source_link.click()
            popup = popup_event.value
            popup.wait_for_load_state("domcontentloaded")
            require(popup.url.split("#")[0] == href.split("#")[0],
                    "SOURCE_POPUP_URL_MISMATCH: observed=" + popup.url)
            require(popup.locator("body").inner_text() == "Offline source target probe",
                    "SOURCE_POPUP_NOT_INTERCEPTED")
            require(any(url.split("#")[0] == href.split("#")[0]
                        for url in source_requests), "SOURCE_REQUEST_NOT_OBSERVED")
            receipt["checks"].append("frozen_source_link_real_popup_no_network")

            require(not page_errors, "BROWSER_SCRIPT_ERROR: " + repr(page_errors))
            require(not external_requests,
                    "UNEXPECTED_NETWORK_ACCESS: " + repr(external_requests))
            receipt["checks"].append("no_script_errors_or_external_requests")
            receipt.update(status="PASS", source_commit=data["source_commit"],
                           indexed_symbols=len(nodes),
                           reviewed_symbols=data["coverage"]["source_reviewed_labels"],
                           relationship_count=len(data["relationships"]),
                           traversed_edge_kind=origin_edge["kind"],
                           source_popup_network="intercepted_offline")
            context.close()
            browser.close()
    except Exception as exc:
        receipt["failure"] = str(exc)
        if page is not None:
            try:
                screenshot = args.receipt.with_suffix(".png")
                page.screenshot(path=str(screenshot), full_page=True, timeout=5000)
                receipt["failure_screenshot"] = str(screenshot)
            except Exception:
                pass
        raise
    finally:
        args.receipt.write_text(json.dumps(receipt, indent=2) + "\n",
                                encoding="utf-8")
        print("WINGLESS_ATLAS_BROWSER=" + receipt["status"])
        print("WINGLESS_ATLAS_BROWSER_RECEIPT=" + str(args.receipt))


if __name__ == "__main__":
    main()

#!/usr/bin/env python3
"""Headless Chrome UI acceptance of source-bound Wingless offline ICE atlases.

Invoked by the existing Windows atlas qualification in its disposable PR checkout.
No application process, inference model, science files, or service is started.
"""
import argparse
import json
import os
import subprocess
import tempfile
from pathlib import Path

from selenium import webdriver
from selenium.common.exceptions import TimeoutException
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.select import Select
from selenium.webdriver.support.ui import WebDriverWait

MERGED_MAIN = "77f7bcd8d72672d6dc0b9427cd0e38e1a92d72b2"
FROZEN_SOURCE = "b87899b50d2c1558e26d74a1921623786fde1282"


def require(value, reason):
    if not value:
        raise AssertionError("ATLAS_CHROMIUM_FAIL_CLOSED: " + reason)


def git(root, *args):
    return subprocess.check_output(
        ["git", *args], cwd=root, text=True, stderr=subprocess.PIPE
    ).strip()


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--root", type=Path, required=True)
    opts = ap.parse_args()
    root = opts.root.resolve()
    head = git(root, "rev-parse", "HEAD")
    require(git(root, "cat-file", "-t", MERGED_MAIN) == "commit",
            "merged main object unavailable")
    require(git(root, "merge-base", MERGED_MAIN, head) == MERGED_MAIN,
            "browser acceptance not based on the merged main")
    # The frozen merged baseline is immutable. A reviewed source-bound
    # regeneration of the PR atlas is permitted, never a silent mutation of
    # the frozen evidence or an unpinned HTML replacement.
    for file in (".ice/code-intelligence.html", ".ice/code-intelligence.json"):
        require(git(root, "cat-file", "-t", f"{MERGED_MAIN}:{file}") == "blob",
                "original promoted atlas unavailable: " + file)
        require(git(root, "cat-file", "-t", f"{head}:{file}") == "blob",
                "candidate atlas unavailable: " + file)
    atlas_file = root / ".ice/code-intelligence.html"
    require(atlas_file.is_file(), "offline atlas unavailable")
    atlas = json.loads((root / ".ice/code-intelligence.json").read_text("utf-8"))
    source_commit = atlas.get("source_commit", "")
    require(git(root, "cat-file", "-t", source_commit) == "commit" and
            git(root, "merge-base", FROZEN_SOURCE, source_commit) == FROZEN_SOURCE and
            git(root, "merge-base", source_commit, head) == source_commit and
            len(atlas["nodes"]) == 352,
            "wrong atlas source lineage or symbol inventory")
    source_header = "SOURCE_COMMIT: " + source_commit
    require((root / ".ice/code-intelligence.ice").read_text("utf-8").splitlines().count(source_header) == 1,
            "candidate atlas narrative source identity mismatch")

    paths = [
        Path(r"C:\Program Files\Google\Chrome\Application\chrome.exe"),
        Path(r"C:\Program Files (x86)\Google\Chrome\Application\chrome.exe"),
    ]
    binary = next((p for p in paths if p.is_file()), None)
    require(os.name == "nt" and binary is not None,
            "Google Chrome missing from hosted Windows environment")

    recorded = ["frozen_baseline_preserved_and_source_bound_candidate_identity"]
    with tempfile.TemporaryDirectory(prefix="wingless-atlas-chrome-") as profile:
        chrome_options = Options()
        chrome_options.binary_location = str(binary)
        for flag in ("--headless=new", "--disable-gpu", "--no-first-run",
                     "--no-default-browser-check", "--disable-extensions",
                     "--disable-background-networking", "--window-size=1440,1000",
                     "--user-data-dir=" + profile):
            chrome_options.add_argument(flag)

        with webdriver.Chrome(options=chrome_options) as browser:
            browser.set_page_load_timeout(25)
            wait = WebDriverWait(browser, 18)
            browser.get(atlas_file.as_uri())
            wait.until(lambda d: d.find_element(By.ID, "nodesCount").text == "352")
            require(browser.find_element(By.ID, "linksCount").text == "155",
                    "relationship total did not render")
            require("352 of 352" in browser.find_element(By.ID, "shown").text,
                    "initial unfiltered list not rendered")
            recorded.append("headless_chromium_offline_boot_and_full_index")

            q = browser.find_element(By.ID, "q")
            q.click()
            q.send_keys("SubmitWork")
            wait.until(lambda d: "4 of 352" in d.find_element(By.ID, "shown").text)
            matches = browser.find_elements(By.CSS_SELECTOR, "button.symbol")
            require(any("Admit and start one bounded inference job" in b.text
                        for b in matches), "typed search failed to find SubmitWork")
            recorded.append("keyboard_search_and_matching_results")

            Select(browser.find_element(By.ID, "basis")).select_by_value("source-reviewed")
            Select(browser.find_element(By.ID, "kind")).select_by_value("method")
            wait.until(lambda d: "1 of 352" in d.find_element(By.ID, "shown").text)
            matching = browser.find_elements(By.CSS_SELECTOR, "button.symbol")
            require(len(matching) == 1 and
                    "Admit and start one bounded inference job" in matching[0].text,
                    "evidence/kind combined filter is not selective")
            recorded.append("two_combined_ui_filters")

            matching[0].click()
            wait.until(lambda d: "Admit and start one bounded inference job" ==
                       d.find_element(By.CSS_SELECTOR, "#detail h3").text)
            worker_link = browser.find_element(By.CSS_SELECTOR, "#detail a").get_attribute("href")
            require(f"/blob/{FROZEN_SOURCE}/worker/service.go#L57-L95" in worker_link,
                    "selected symbol has a different source definition link")
            recorded.append("actual_symbol_selection_and_source_anchor")

            related = browser.find_elements(By.CSS_SELECTOR, "#detail .related button")
            target = next((x for x in related
                           if "Execute bounded candidate inference" in x.text), None)
            require(target is not None, "typed SubmitWork-to-Runner.Run edge missing")
            target.click()
            wait.until(lambda d: d.find_element(By.CSS_SELECTOR, "#detail h3").text ==
                       "Execute bounded candidate inference with deterministic verification")
            source = browser.find_element(By.CSS_SELECTOR, "#detail a")
            expected_link = source.get_attribute("href")
            require(f"/blob/{FROZEN_SOURCE}/broker/runner.go#L33-L112" in expected_link,
                    "cross-file graph traversal led to wrong pinned source")
            recorded.append("actual_click_cross_file_relationship")

            current_windows = set(browser.window_handles)
            source.click()
            wait.until(lambda d: len(set(d.window_handles) - current_windows) == 1)
            next_tab = next(iter(set(browser.window_handles) - current_windows))
            browser.switch_to.window(next_tab)
            wait.until(lambda d: d.current_url.startswith(expected_link.split("#")[0]))
            require(browser.current_url.startswith("https://github.com/"),
                    "source link did not actually open GitHub in a second tab")
            recorded.append("actual_click_source_link_new_tab")

            browser.switch_to.window(next(iter(current_windows)))
            if os.getenv("RUNNER_TEMP"):
                browser.save_screenshot(
                    str(Path(os.environ["RUNNER_TEMP"]) / "wingless-atlas-chrome.png"))

    require(git(root, "status", "--porcelain") == "",
            "browser tests changed the repository checkout")
    print(json.dumps({
        "result": "WINGLESS_ATLAS_CHROMIUM_BROWSER_ACCEPTANCE_PASS",
        "merged_main": MERGED_MAIN,
        "browser_test_commit": head,
        "source_commit": source_commit,
        "tested_html": ".ice/code-intelligence.html",
        "browser": "Google Chrome, Windows hosted, real Selenium UI",
        "steps": recorded,
        "original_checkout_dirty": False,
        "network_requirement": "only optional GitHub source-link navigation",
    }, sort_keys=True))


if __name__ == "__main__":
    main()

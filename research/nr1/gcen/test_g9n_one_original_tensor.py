"""G9N original expert selection negative controls, no model-weight downloads."""
import unittest
from unittest.mock import patch
from g9n_one_original_tensor import choose_target, MAX_TENSOR_BYTES


class OriginalExpertSelectionTests(unittest.TestCase):
    def make_header(self):
        name="model.layers.0.mlp.experts.0.down_proj.weight"
        return name,{name:{"dtype":"BF16","shape":[128,256]}},{name:"model-00001-of-00016.safetensors"}

    def test_frozen_source_tensor_selection(self):
        key,header,index=self.make_header()
        with patch("g9n_one_original_tensor.select_experts",
                   return_value={"names":[key],"layer":0}):
            name,shape,layer=choose_target(header,index)
            self.assertEqual((name,shape,layer),(key,(128,256),0))

    def test_reject_non_bf16(self):
        key,header,index=self.make_header()
        header[key]["dtype"]="F32"
        with patch("g9n_one_original_tensor.select_experts",
                   return_value={"names":[key],"layer":0}):
            with self.assertRaisesRegex(ValueError,"BF16"):
                choose_target(header,index)

    def test_reject_oversized_original_tensor(self):
        key,header,index=self.make_header()
        header[key]["shape"]=[MAX_TENSOR_BYTES,2]
        with patch("g9n_one_original_tensor.select_experts",
                   return_value={"names":[key],"layer":0}):
            with self.assertRaisesRegex(ValueError,"bounds"):
                choose_target(header,index)

    def test_reject_cross_shard_reference(self):
        key,header,index=self.make_header()
        index[key]="model-00002-of-00016.safetensors"
        with patch("g9n_one_original_tensor.select_experts",
                   return_value={"names":[key],"layer":0}):
            with self.assertRaisesRegex(ValueError,"bounds"):
                choose_target(header,index)

    def test_reject_missing_tensor_metadata(self):
        key,header,index=self.make_header()
        header={}
        with patch("g9n_one_original_tensor.select_experts",
                   return_value={"names":[key],"layer":0}):
            with self.assertRaisesRegex(ValueError,"BF16"):
                choose_target(header,index)


if __name__ == "__main__":
    unittest.main()

// Mock for remark-gfm plugin
export default function remarkGfm() {
  return function (tree) {
    // Mock implementation - just return the tree as-is
    return tree;
  };
}

import { animate, τmax } from "../motion.js";

export default function (el) {
  const tree = el.querySelector("g");
  animate(tree, { rotate: [0, -10, 6, -3, 0] }, { duration: τmax });
}

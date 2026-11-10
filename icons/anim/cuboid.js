import { animate, τmax } from "../motion.js";

export default function (el) {
  const edges = el.querySelector("g");
  animate(edges, { scale: [0, 1] }, { duration: τmax });
}

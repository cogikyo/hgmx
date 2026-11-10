import { animate, τmax } from "../motion.js";

export default function (el) {
  const hop = el.querySelector("g");
  animate(hop, { rotate: [0, 8, -5, 2, 0] }, { duration: τmax });
}

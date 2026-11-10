import { animate, τmax } from "../motion.js";

export default function (el) {
  const torch = el.querySelector("g");
  animate(torch, { rotate: [0, -12, 8, 0] }, { duration: τmax });
}

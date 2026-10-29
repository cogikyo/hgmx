import { animate, τmax } from "../motion.js";

export default function (el) {
  const tag = el.querySelector("g");
  animate(tag, { rotate: [0, -18, 14, -8, 0] }, { duration: τmax });
}

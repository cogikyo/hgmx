import { animate, τmax } from "../motion.js";

export default function (el) {
  const anchor = el.querySelector("g");
  animate(anchor, { rotate: [0, -10, 7, -3, 0] }, { duration: τmax });
}

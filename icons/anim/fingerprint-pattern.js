import { animate, τmax } from "../motion.js";

export default function (el) {
  const print = el.querySelector("g");
  animate(print, { scale: [1, 0.9, 1.05, 1] }, { duration: τmax });
}

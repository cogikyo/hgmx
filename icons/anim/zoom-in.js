import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const lens = el.querySelector("g");
  animate(lens, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });

  const plus = el.querySelectorAll(":scope > line");
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}

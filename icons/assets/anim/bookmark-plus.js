import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const ribbon = el.querySelector("g");
  animate(ribbon, { scaleY: [1, 1.15, 0.95, 1] }, { duration: τmax });

  const plus = el.querySelectorAll(":scope > path");
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}

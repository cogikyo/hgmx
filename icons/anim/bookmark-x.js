import { animate, stagger, τmid, τmin, τmax } from "../motion.js";

export default function (el) {
  const ribbon = el.querySelector("g");
  animate(ribbon, { scaleY: [1, 1.15, 0.95, 1] }, { duration: τmax });

  const x = el.querySelectorAll(":scope > path");
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(τmin) });
}

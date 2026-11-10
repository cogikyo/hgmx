import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const lens = el.querySelector("g");
  animate(lens, { scale: [1, 1.15, 0.85, 1] }, { duration: τmax });

  const minus = el.querySelector(":scope > line");
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}

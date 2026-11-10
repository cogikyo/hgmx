import { animate, τmid, τmin, τmax } from "../motion.js";

export default function (el) {
  const ribbon = el.querySelector("g");
  animate(ribbon, { scaleY: [1, 1.15, 0.95, 1] }, { duration: τmax });

  const minus = el.querySelector(":scope > path");
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}

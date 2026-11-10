import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const clock = el.querySelector("g");
  animate(clock, { rotate: [0, -12, 12, -8, 8, 0] }, { duration: τmax });

  const minus = el.querySelector('path[d="M9 13h6"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}

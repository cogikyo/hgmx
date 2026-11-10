import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const map = el.querySelector("g");
  animate(map, { scaleX: [0.25, 1.15, 1] }, { duration: τmid });

  const minus = el.querySelector('path[d="M21 18h-6"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}

import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const wall = el.querySelector("g");
  animate(wall, { scaleY: [0, 1.1, 1] }, { duration: τmid });

  const flame = el.querySelector('path[d^="M17 9"]');
  animate(flame, { scaleY: [1, 0.85, 1.15, 1] }, { duration: τmid, delay: τmin });
}

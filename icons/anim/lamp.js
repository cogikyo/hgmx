import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const base = el.querySelector('path[d^="M8 20"]');
  animate(base, { scaleY: [1, 0.7, 1] }, { duration: τmin });

  const shade = el.querySelector('path[d^="M4.077"]');
  animate(shade, { scale: [1, 0.85, 1.15, 1] }, { duration: τmid, delay: τmin });
}

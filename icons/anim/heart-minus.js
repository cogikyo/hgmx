import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const heart = el.querySelector('path[d^="m14.876"]');
  animate(heart, { scale: [1, 1.25, 0.85, 1] }, { duration: τmid });

  const minus = el.querySelector('path[d="M15 15h6"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}

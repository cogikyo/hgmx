import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const heart = el.querySelector('path[d^="m14.479"]');
  animate(heart, { scale: [1, 1.25, 0.85, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M15 15h6"], path[d="M18 12v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}

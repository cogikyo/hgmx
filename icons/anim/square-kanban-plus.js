import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M12.354"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M15 18h6"], path[d="M18 21v-6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}

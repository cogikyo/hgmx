import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M9 10h6"], path[d="M12 13V7"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });

  const minus = el.querySelector('path[d="M9 17h6"]');
  animate(minus, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}

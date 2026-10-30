import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d^="M9 15"], path[d^="M12 18"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}

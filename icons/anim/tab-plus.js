import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const tab = el.querySelector('path[d^="M4 20"]');
  animate(tab, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M15 12H9"], path[d="M12 9v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}

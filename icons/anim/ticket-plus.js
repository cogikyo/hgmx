import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ticket = el.querySelector('path[d^="M2 9"]');
  animate(ticket, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M9 12h6"], path[d="M12 9v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}

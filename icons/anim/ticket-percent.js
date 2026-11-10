import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ticket = el.querySelector('path[d^="M2 9"]');
  animate(ticket, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d="m15 9-6 6"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const dots = [el.querySelector('path[d="M9 9h.01"]'), el.querySelector('path[d="M15 15h.01"]')];
  animate(dots, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}

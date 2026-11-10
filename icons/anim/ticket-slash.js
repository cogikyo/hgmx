import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ticket = el.querySelector('path[d^="M2 9"]');
  animate(ticket, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d="m9.5 14.5 5-5"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

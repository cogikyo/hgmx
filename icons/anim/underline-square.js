import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const subject = el.querySelectorAll('rect, path[d^="M15 7"]');
  animate(subject, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const rule = el.querySelector('path[d="M7 17h10"]');
  animate(rule, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const subject = el.querySelectorAll('rect, path[d^="M10.5 12"], path[d^="M8 17"]');
  animate(subject, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const strike = el.querySelector('path[d="M7 12h10"]');
  animate(strike, { scaleX: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

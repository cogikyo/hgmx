import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector("rect");
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const rule = el.querySelector('path[d="M3 9h18"]');
  animate(rule, { scaleX: [0, 1] }, { duration: τmid, delay: τmin });
}

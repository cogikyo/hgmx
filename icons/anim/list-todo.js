import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="h8"]');
  animate(lines, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const box = el.querySelector("rect");
  animate(box, { scale: [0, 1.25, 1] }, { duration: τmid });
  const tick = el.querySelector('path[d^="m3 17"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

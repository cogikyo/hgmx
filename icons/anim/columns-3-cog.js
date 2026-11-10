import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M10.6 "]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const splits = [el.querySelector('path[d="M9 3v18"]'), el.querySelector('path[d="M15 3v7.6"]')];
  animate(splits, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });

  const gear = el.querySelector("g.gear");
  animate(gear, { rotate: [0, 360], scale: [1, 1.25, 1] }, { duration: τmid, delay: τmin });
}

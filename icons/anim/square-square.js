import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('rect[x="3"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const inner = el.querySelector('rect[x="8"]');
  animate(inner, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

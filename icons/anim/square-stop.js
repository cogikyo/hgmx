import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('rect[x="3"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stop = el.querySelector('rect[x="9"]');
  animate(stop, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}

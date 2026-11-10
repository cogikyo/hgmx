import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector('rect[x="2"]');
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stop = el.querySelector('rect[x="9"]');
  animate(stop, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}

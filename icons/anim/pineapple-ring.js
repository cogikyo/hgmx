import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelectorAll('ellipse, path[d^="M2 10"]');
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const ticks = el.querySelectorAll('path[d^="m6 4"], path[d^="M16.7"], path[d^="m8 15"], path[d^="M12 16"]');
  animate(ticks, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.04, { startDelay: τmid }) });
}

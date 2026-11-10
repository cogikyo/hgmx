import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const grid = el.querySelector('path[d^="M12 3"]');
  animate(grid, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d^="m16.5 21.5"]'), el.querySelector('path[d^="m16.5 16.5"]')];
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(0.15) });
}

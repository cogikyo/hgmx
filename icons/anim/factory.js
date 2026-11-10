import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M3 19"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const windows = ["M8", "M12", "M16"].map((x) => el.querySelector(`path[d="${x} 16h.01"]`));
  animate(windows, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}

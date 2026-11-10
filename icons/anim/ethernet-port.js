import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const port = el.querySelector('path[d^="M19 17"]');
  animate(port, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pins = ["M6", "M10", "M14", "M18"].map((x) => el.querySelector(`path[d="${x} 8v1"]`));
  animate(pins, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.04, { startDelay: τmin }) });
}

import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const canopy = el.querySelector('path[d^="M2.4 14.4"]');
  animate(canopy, { scaleX: [0.3, 1.15, 1] }, { duration: τmid });

  const sun = el.querySelector("circle");
  animate(sun, { scale: [1, 1.35, 1] }, { duration: τmid, delay: τmin });
}

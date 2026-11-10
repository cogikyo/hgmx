import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const meridian = el.querySelector('path[d^="M12 2"]');
  animate(meridian, { scaleX: [1, -1, 1] }, { duration: τmax });
}

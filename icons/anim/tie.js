import { animate, τmax } from "../motion.js";

export default function (el) {
  const band = el.querySelectorAll('path[d^="M8.85"], path[d^="m9.7"]');
  animate(band, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}

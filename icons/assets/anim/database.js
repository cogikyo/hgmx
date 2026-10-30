import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M3 5V19"]');
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const band = el.querySelector('path[d^="M3 12"]');
  animate(band, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const lid = el.querySelector("ellipse");
  animate(lid, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });
}

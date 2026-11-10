import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll('path[d^="M3 "], path[d^="M21 5"]');
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector("ellipse");
  animate(lid, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });

  const bolt = el.querySelector('path[d^="M21 12"]');
  animate(bolt, { pathLength: [0, 1], opacity: [0, 1, 0.25, 1] }, { duration: τmid, delay: τmin });
}

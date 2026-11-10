import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll('path[d^="M3 "], path[d^="M21"]');
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector("ellipse");
  animate(lid, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });

  const arrow = el.querySelectorAll('path[d^="m16 19"], path[d^="M19 16"]');
  animate(arrow, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}

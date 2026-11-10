import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll('path[d^="M3 "], path[d^="M21"]');
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector("ellipse");
  animate(lid, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });

  const plus = el.querySelectorAll('path[d^="M19 16"], path[d^="M22 19"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}

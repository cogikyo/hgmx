import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll('path[d^="M3 "], path[d^="M21"]');
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector("ellipse");
  animate(lid, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });

  const arc = el.querySelector('path[d^="M13 20"]');
  animate(arc, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const head = el.querySelector('path[d^="M12 12"]');
  animate(head, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}

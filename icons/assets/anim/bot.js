import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelector("rect");
  animate(head, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const antenna = el.querySelector('path[d^="M12 8V4"]');
  animate(antenna, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const eyes = el.querySelectorAll('path[d$="v2"]');
  animate(eyes, { scaleY: [1, 0.1, 1] }, { duration: τmin, delay: τmid });
}

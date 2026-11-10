import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const sticker = el.querySelectorAll('path[d^="M21 9"], path[d^="M15 3"]');
  animate(sticker, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const mouth = el.querySelector('path[d^="M10 16"]');
  animate(mouth, { scaleY: [0, 1.15, 1] }, { duration: τmid, delay: τmin });

  const eyes = el.querySelectorAll('path[d$="13h.01"]');
  animate(eyes, { translateY: [1, -1, 0] }, { duration: τmin, delay: τmid });
}

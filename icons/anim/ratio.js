import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const portrait = el.querySelector('rect[width="12"]');
  animate(portrait, { scaleY: [1, 0.85, 1] }, { duration: τmid });

  const landscape = el.querySelector('rect[width="20"]');
  animate(landscape, { scaleX: [1, 0.85, 1] }, { duration: τmid, delay: τmin });
}

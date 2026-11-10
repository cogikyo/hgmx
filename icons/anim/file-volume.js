import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M4 11.55"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const speaker = el.querySelector('path[d^="M8 14.502"]');
  animate(speaker, { scale: [1, 0.85, 1] }, { duration: τmid });

  const wave = el.querySelector('path[d^="M12 15"]');
  animate(wave, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}

import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const crown = el.querySelector('path[d^="M11.562"]');
  animate(crown, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: τmax });
  const band = el.querySelector('path[d="M5 21h14"]');
  animate(band, { scaleX: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}

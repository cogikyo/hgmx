import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M4 6.3"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const band = el.querySelector('path[d^="M12 16"]');
  animate(band, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const lid = el.querySelector("ellipse");
  animate(lid, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });
}

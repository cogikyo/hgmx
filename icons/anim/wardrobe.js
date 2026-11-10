import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const wardrobe = el.querySelector("rect");
  animate(wardrobe, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const seam = el.querySelector("path[transform-origin]");
  animate(seam, { scaleY: [0, 1] }, { duration: τmid, delay: τmin });
}

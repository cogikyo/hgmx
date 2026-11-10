import { animate, stagger } from "../motion.js";

export default function (el) {
  const at = (x, y) => el.querySelector(`circle[cx="${x}"][cy="${y}"]`);
  const tiles = [at(6.5, 6.5), at(17.5, 6.5), at(17.5, 17.5), at(6.5, 17.5)];
  animate(tiles, { pathLength: [0, 1], opacity: [0, 1], rotate: [270, 90] }, { duration: 0.3, delay: stagger(0.06) });
}
